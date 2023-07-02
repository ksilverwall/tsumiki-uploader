import { BackendApi } from "./BackendApi";
import { S3SignedAccessor } from "./S3SignedAccessor";
import { BufferUploader } from "./Uploader";


export default class BackendSignedBufferUploader implements BufferUploader {
  constructor() { }
  async upload(zipData: ArrayBuffer): Promise<string> {
    const endPoint = import.meta.env.VITE_BACKEND_API_ENDPOINT;
    const backendApi = new BackendApi(endPoint);
    const s3 = new S3SignedAccessor();
    const t = await backendApi.createTransaction();
    await s3.put(t.url, zipData);

    return `${t.id}`;
  }
}
