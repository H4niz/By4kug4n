import click
import yaml
from pathlib import Path

@click.group()
@click.option('--config', type=click.Path(exists=True), help='Config file path')
@click.pass_context
def cli(ctx, config):
    """Byakugan - Advanced API Security Scanner"""
    ctx.ensure_object(dict)
    
    if config:
        with open(config) as f:
            ctx.obj['config'] = yaml.safe_load(f)
    else:
        # Load default config
        default_config = Path.home() / '.byakugan' / 'config.yaml'
        if default_config.exists():
            with open(default_config) as f:
                ctx.obj['config'] = yaml.safe_load(f)