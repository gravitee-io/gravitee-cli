import logging

import click

from graviteeio_cli.exeptions import GraviteeioError
from graviteeio_cli.graviteeio.modules import GraviteeioModule
from graviteeio_cli.graviteeio.output import OutputFormatType
from graviteeio_cli.graviteeio.profiles import Auth_Type


@click.command()
@click.option('--user', help='authentication user', required=True)
@click.option('--type', 'auth_type',
                        help='authentication type', required=True,
                    type=click.Choice(Auth_Type.list_name(), case_sensitive=False))
@click.pass_obj
def create(obj, user, auth_type):
    """This command create a new authentication user"""
    config = obj['config'].getGraviteeioConfig(GraviteeioModule.APIM)
    auth_list = config.get_auth_list()

    for auth in auth_list:
        if auth["username"] == user and auth_type.upper() == auth["type"].upper():
            raise GraviteeioError("Username [{}] already exit for authentication type {}.".format(user, auth_type))
    
    auth_list.append({
        "username": user,
        "type": auth_type,
        "is_active": False
    })

    config.save(auth = auth_list)

    click.echo("User [{}] saved.".format(user))
