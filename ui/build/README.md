# Build Resources

## Icon Files

- `appicon.svg` - Application icon (SVG format)
- `trayicon.png` - System tray icon (PNG format, recommended 22x22 or 16x16)

### Generate Tray Icon

If you need to generate a tray icon, you can use the following commands to convert SVG to PNG:

```bash
# Using ImageMagick
convert -background none -resize 22x22 appicon.svg trayicon.png

# Or using rsvg-convert
rsvg-convert -w 22 -h 22 appicon.svg > trayicon.png
```

### Notes

- Windows tray icons should preferably use ICO format
- macOS and Linux support PNG format
- Recommended icon size is 22x22 pixels

