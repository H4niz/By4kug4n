import click
from ..client.core import CoreClient

@click.command()
@click.argument('target')
@click.option('--rules', '-r', multiple=True, help='Security rules to apply')
@click.pass_context
def scan(ctx, target, rules):
    """Execute security scan on target"""
    config = ctx.obj.get('config', {})
    core_addr = config.get('core', {}).get('address', 'localhost:50051')

    client = CoreClient(core_addr)
    try:
        result = client.execute_scan(target, rules)
        for finding in result.findings:
            click.echo(f"[{finding.severity}] {finding.description}")
    finally:
        client.close()