from pathlib import Path
from PIL import Image

for p in Path('.').glob('*.bmp'):
    with Image.open(p) as im:
        im.save(p.with_suffix('.tif'), compression='tiff_lzw')