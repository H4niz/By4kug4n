import click
from byakugan.cmd.cli.commands.root import cli
from byakugan.cmd.cli.commands.scan import scan

cli.add_command(scan)

if __name__ == '__main__':
    cli(obj={})