# OpenAI Supported File Types for Vector Stores

OpenAI's vector stores support a variety of file types for efficient data retrieval and processing. The supported file formats include:

## Text and Code Files

- `.c`
- `.cpp`
- `.css`
- `.csv`
- `.go`
- `.html`
- `.java`
- `.js`
- `.json`
- `.md`
- `.php`
- `.py`
- `.rb`
- `.sh`
- `.ts`
- `.txt`
- `.xml`
- `.tex`

## Document Files

- `.doc`
- `.docx`
- `.pdf`
- `.pptx`
- `.xlsx`

## Image Files

- `.gif`
- `.jpeg`
- `.jpg`
- `.png`
- `.webp`

## Compressed Files

- `.tar`
- `.zip`

**Note:** While these formats are supported, certain file types, such as `.csv` and `.xlsx`, have specific considerations. Users have reported encountering errors when uploading these formats directly to vector stores, despite official documentation indicating support. For instance, attempts to upload `.csv` and `.xlsx` files have sometimes resulted in unsupported file format errors. [Source](https://community.openai.com/t/what-file-types-are-actually-supported/929529)

To mitigate such issues, consider the following recommendations:

- **Conversion:** Convert files to supported formats like `.txt` or `.json` before uploading.
- **Code Interpreter:** Enable the code interpreter feature to facilitate the processing of these file types.

Always ensure that your files are encoded in `utf-8`, `utf-16`, or ASCII to prevent encoding-related issues during the upload process.
