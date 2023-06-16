import JSZip from "jszip";
import { uuidv7 } from "uuidv7";

class ZipUploader {
  async upload(formData: FormData): Promise<string> {
    console.log("dummy impl of Uploader::upload()");
    return uuidv7();
  }
}

class Uploader {
  async upload(files: File[]) {
    const zip = new JSZip();

    files.forEach((file) => {
      zip.file(file.name, file);
    });

    const content = await zip.generateAsync({ type: "blob" });

    const formData = new FormData();
    formData.append("file", content, "archive.zip");

    return new ZipUploader().upload(formData);
  }
}

export default Uploader;
