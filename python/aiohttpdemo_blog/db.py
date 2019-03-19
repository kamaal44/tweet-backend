import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

__all__ = ['tweet', 'profile', 'tag', 'comment', 'friend', 'follower']

meta = sa.MetaData()

tweet = sa.Table(
    'tweet', meta,
    sa.Column('id', sa.Integer, nullable=False),
    sa.Column('profile_id', sa.Integer, nullable=False),
    sa.Column('title', sa.String(200), nullable=False),
    sa.Column('teaser', sa.String(500), nullable=False),
    sa.Column('body', sa.Text, nullable=False),
    sa.Column('views', sa.Integer, nullable=False),
    sa.Column('average_note', sa.Float, nullable=False),
    sa.Column('pictures', postgresql.JSON, server_default='{}'),
    sa.Column('published_at', sa.Date, nullable=False),
    sa.Column('tags', postgresql.ARRAY(sa.Integer), server_default='{}'),
    sa.Column('category', sa.String(50), nullable=False),
    sa.Column('subcategory', sa.String(50), nullable=False),
    sa.Column('backlinks', postgresql.JSON, server_default='{}'),
    # Indexes #
    sa.PrimaryKeyConstraint('id', name='tweet_id_pkey')
)

profile = sa.Table(
    'profile', meta,
    sa.Column('id', sa.Integer, nullable=False),
    sa.Column('name', sa.String(200), nullable=False),
    # sa.Column('pictures', postgresql.JSON, server_default='{}'),
    # sa.Column('published_at', sa.Date, nullable=False),
    # sa.Column('tags', postgresql.ARRAY(sa.Integer), server_default='{}'),
    sa.Column('active', sa.Boolean, nullable=False,
              server_default='FALSE'),
    # Indexes #
    sa.PrimaryKeyConstraint('id', name='profile_id_pkey')
)

tag = sa.Table(
    'tag', meta,
    sa.Column('id', sa.Integer, nullable=False),
    sa.Column('name', sa.String(10), nullable=False),
    sa.Column('published', sa.Boolean, nullable=False,
              server_default='FALSE'),

    # Indexes #
    sa.PrimaryKeyConstraint('id', name='tag_id_pkey')
)

comment = sa.Table(
    'comment', meta,
    sa.Column('id', sa.Integer, nullable=False),
    sa.Column('tweet_id', sa.Integer, nullable=False),
    sa.Column('author', postgresql.JSON, server_default='{}'),
    sa.Column('body', sa.Text, nullable=False),
    sa.Column('created_at', sa.Date, nullable=False),

    # Indexes #
    sa.PrimaryKeyConstraint('id', name='comment_id_pkey'),
    sa.ForeignKeyConstraint(['tweet_id'], [tweet.c.id], name='tweet_fkey',
                            ondelete='CASCADE')
)

friend = sa.Table(
    'friend', meta,
    sa.Column('id', sa.Integer, nullable=False),
    sa.Column('profile_id', sa.Integer, nullable=False),
    sa.Column('friend_id', sa.Integer, nullable=False),
    # Indexes #
    sa.PrimaryKeyConstraint('id', name='friend_id_pkey')
)

follower = sa.Table(
    'follower', meta,
    sa.Column('id', sa.Integer, nullable=False),
    sa.Column('follower_id', sa.Integer, nullable=False),

    # Indexes #
    sa.PrimaryKeyConstraint('id', name='follower_id_pkey')
)


