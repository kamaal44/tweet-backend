import asyncio
import pathlib

import psycopg2
from faker import Factory
from sqlalchemy.schema import CreateTable, DropTable

import db
from utils import init_postgres


PROJ_ROOT = pathlib.Path(__file__).parent.parent

conf = {"host": "127.0.0.1",
        "port": 8080,
        "postgres": {
            "database": "aiohttp_admin",
            "user": "aiohttp_admin_user",
            "password": "mysecretpassword",
            "host": "127.0.0.1",
            "port": 5432,
            "minsize": 1,
            "maxsize": 5}
        }


async def delete_tables(pg, tables):
    async with pg.acquire() as conn:
        for table in reversed(tables):
            drop_expr = DropTable(table)
            try:
                await conn.execute(drop_expr)
            except psycopg2.ProgrammingError:
                pass


async def prepare_tables(pg):
    tables = [db.tweet, db.profile, db.tag, db.comment]
    await delete_tables(pg, tables)
    async with pg.acquire() as conn:
        for table in tables:
            create_expr = CreateTable(table)
            await conn.execute(create_expr)


async def insert_data(pg, table, values):
    async with pg.acquire() as conn:
        query = table.insert().values(values).returning(table.c.id)
        cursor = await conn.execute(query)
        resp = await cursor.fetchall()
    return [r[0] for r in resp]


async def generate_profiles(pg, rows, fake):
    values = []
    for i in range(rows):
        values.append({
            'name': fake.word()[:10],
            'active': bool(i % 2),
        })
    print(values)
    ids = await insert_data(pg, db.profile, values)
    return ids

async def generate_tags(pg, rows, fake):
    values = []
    for i in range(rows):
        values.append({
            'name': fake.word()[:10],
            'published': bool(i % 2),
        })
    ids = await insert_data(pg, db.tag, values)
    return ids


async def generate_tweets(pg, rows, fake, tag_ids, profile_ids):
    values = []
    for i in range(rows):
        values.append({
            'profile_id': 1, # [profile_ids[(i + j) % len(profile_ids)] for j in range(1)],
            'title': fake.sentence(nb_words=7)[:100],
            'teaser': fake.paragraph(nb_sentences=4)[:100],
            'body': fake.text(max_nb_chars=280),
            'views': i % 1000,
            'average_note': i % 0.1,
            'pictures': {'first': {'name': fake.word(),
                                   'url': fake.image_url()}},
            'published_at': fake.iso8601(),
            'tags': [tag_ids[(i + j) % len(tag_ids)] for j in range(7)],
            'category': fake.word()[:50],
            'subcategory': fake.word()[:50],
            'backlinks': {'date': fake.iso8601(),
                          'url': fake.uri()},
        })
    print(values)
    ids = await insert_data(pg, db.tweet, values)
    return ids


async def generate_comments(pg, rows, fake, tweet_ids):
    values = []
    for tweet_id in tweet_ids:
        for i in range(rows):
            values.append({
                'tweet_id': tweet_id,
                'body': fake.text(max_nb_chars=280),
                'created_at': fake.iso8601(),
                'author': {'name': fake.name(),
                           'email': fake.email()},
            })

    await insert_data(pg, db.comment, values)


async def init(loop):
    print("Generating Fake Data")
    pg = await init_postgres(conf['postgres'], loop)
    fake = Factory.create()
    fake.seed(1234)
    await prepare_tables(pg)

    rows = 1000

    tag_ids = await generate_tags(pg, 500, fake)
    profile_ids = await generate_profiles(pg, 500, fake)
    tweet_ids = await generate_tweets(pg, rows, fake, tag_ids, profile_ids)
    await generate_comments(pg, 25, fake, tweet_ids)


def main():
    loop = asyncio.get_event_loop()
    loop.run_until_complete(init(loop))


if __name__ == "__main__":
    main()
