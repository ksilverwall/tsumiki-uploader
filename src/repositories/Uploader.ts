import JSZip from "jszip";

export interface BufferUploader {
  upload(zipData: ArrayBuffer): Promise<string>;
}

class Uploader {
  constructor(private readonly uploader: BufferUploader) {}
  async upload(files: File[]) {
    const zip = new JSZip();

    files.forEach((file) => {
      zip.file(file.name, file);
    });

    const zipData = await zip.generateAsync({ type: "arraybuffer" });

    return this.uploader.upload(zipData);
  }
}

export default Uploader;
