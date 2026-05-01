# img2jpg

A fast, lightweight Go CLI utility for converting images to JPG format. Designed for batch processing and seamless integration with AI agents like Gemini CLI.

## Prerequisites
- Go 1.25 or higher
- Linux x86_64 (Xubuntu/Ubuntu tested)

## Installation

### 1. Clone and Build
```bash
git clone https://github.com/yourusername/img2jpg.git
cd img2jpg
go mod tidy
go build -ldflags="-s -w" -o img2jpg .
```

### 2. Add to PATH
```bash
cp img2jpg ~/.local/bin/
chmod +x ~/.local/bin/img2jpg
source ~/.bashrc
```

## Usage
```bash
img2jpg <input_path> [-q quality] [-o output_dir] [-r]
```

| Flag | Description | Default |
|------|-------------|---------|
| `<input_path>` | Target file or directory | Required |
| `-q` | JPG quality (1-100) | 80 |
| `-o` | Output directory | `<input>/converted` |
| `-r` | Process directories recursively | false |

## Gemini CLI Integration

To enable automatic invocation by Gemini CLI, create a custom skill directory in your project or home configuration.

### 1. Create Skill Directory
```bash
mkdir -p .gemini/skills/image_converter
```

### 2. Add Skill Definition
Create a file at `.gemini/skills/image_converter/SKILL.md` and paste the following content exactly:

```markdown
---
name: img_converter
description: Converts images to JPG using the local `img2jpg` utility. Use this when the user wants to optimize, convert, or batch process images.
---

# Knowledge
- You have access to a local Go utility called `img2jpg`.
- **Command:** `img2jpg <input_path> [-q quality] [-o output_dir] [-r]`
- **Recursive Mode:** Use `-r` when the input is a directory.
- **Output Rule:** Unless specified, set the output directory to `<input_path>/converted`.

# Actions
When the user asks to optimize images:
1. Identify if the target is a file or directory.
2. Construct the `img2jpg` command using the shell tool.
3. Default to `-q 80`.
```

### 3. Verify Integration
Launch Gemini CLI in the directory containing `.gemini`. The agent will automatically load the skill and route image conversion requests to `img2jpg` using your defined rules.

## Supported Formats
PNG, WebP, GIF, TIFF, BMP, ICO. Output is always standard JPEG with white background composition for alpha channels. Existing JPG files are automatically skipped.
