from aiohttp_admin.contrib import models
from aiohttp_admin.backends.sa import PGResource

from .main import schema
from ..db import comment


@schema.register
class Friends(models.ModelAdmin):
    fields = ('id', 'profile_id', 'profile_name', 'friend_id', 'friend_name',)

    can_edit = False
    can_create = False
    can_delete = False

    class Meta:
        resource_type = PGResource
        table = comment
