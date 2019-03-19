from aiohttp_admin.contrib import models
from aiohttp_admin.backends.sa import PGResource

from .main import schema
from ..db import tweet


@schema.register
class Tweets(models.ModelAdmin):
    fields = ('id', 'profile_name', 'pictures', 'backlinks', 'subcategory',)
    can_edit = False
    can_create = True
    can_delete = False

    class Meta:
        resource_type = PGResource
        table = tweet
