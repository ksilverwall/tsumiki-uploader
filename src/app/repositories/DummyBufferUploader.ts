import { uuidv7 } from "uuidv7";
import { BufferUploader } from "./Uploader";


export class DummyBufferUploader implements BufferUploader {
  async upload(formData: ArrayBuffer): Promise<string> {
    console.log("dummy impl of Uploader::upload()");
    return uuidv7();
  }
}
