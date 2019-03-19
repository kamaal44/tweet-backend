from aiohttp_admin.contrib import models
from aiohttp_admin.backends.sa import PGResource

from .main import schema
from ..db import profile


@schema.register
class Profiles(models.ModelAdmin):
    fields = ('id', 'name', 'picture', 'active',)

    class Meta:
        resource_type = PGResource
        table = profile
