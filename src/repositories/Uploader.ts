import JSZip from "jszip";
import { ApplicationError } from "./ApplicationError";

export interface BufferUploader {
  upload(zipData: ArrayBuffer): Promise<string>;
}

class Uploader {
  constructor(private readonly uploader: BufferUploader) { }
  async upload(files: File[]) {
    try {
      const zip = new JSZip();

      files.forEach((file) => {
        zip.file(file.name, file);
      });

      const zipData = await zip.generateAsync({ type: "arraybuffer" });

      return this.uploader.upload(zipData);
    } catch (err) {
      if (err instanceof ApplicationError) {
        err.toStrings().forEach(s => {
          console.error(s)
        });
      } else if (err instanceof Error) {
        console.error('unknown error: ', err.message);
      } else {
        console.error('illegal format err obj handled: ', err);
      }
      throw err;
    }
  }
}

export default Uploader;
