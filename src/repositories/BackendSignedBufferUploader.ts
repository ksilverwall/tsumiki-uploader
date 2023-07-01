import { ApplicationError } from "./ApplicationError";
import { BackendApi } from "./BackendApi";
import { S3SignedAccessor } from "./S3SignedAccessor";
import { BufferUploader } from "./Uploader";


export default class BackendSignedBufferUploader implements BufferUploader {
  constructor() {}
  async upload(zipData: ArrayBuffer): Promise<string> {
    try {
      const endPoint = import.meta.env.VITE_BACKEND_API_ENDPOINT;
      const backendApi = new BackendApi(endPoint);
      const s3 = new S3SignedAccessor();
      const t = await backendApi.createTransaction();
      await s3.put(t.url, zipData);
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
    }
  
    return "";
  }
}
