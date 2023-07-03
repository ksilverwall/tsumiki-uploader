import { BackendInterface } from "../interface";
import { BackendApi } from "./BackendApi";
import { S3SignedAccessor } from "./S3SignedAccessor";

type Accessor = {
  backendApi: BackendApi;
  s3Accessor: S3SignedAccessor;
}

export default class Backend implements BackendInterface {
  constructor(private readonly accessor: Accessor) { }
  async upload(zipData: ArrayBuffer): Promise<string> {
    const t = await this.accessor.backendApi.createTransaction();
    await this.accessor.s3Accessor.put(t.url, zipData);

    return t.id;
  }
  async download(key: string): Promise<void> {
    const url = await this.accessor.backendApi.getDownloadUrl(key);
    const link = document.createElement('a');
    link.href = url;
    link.download = "dummy-file-name";
    link.click();
  }
}
