import { BackendInterface, DownloadInfo } from "../interface";
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
  async download(key: string): Promise<DownloadInfo> {
    const info = await this.accessor.backendApi.getDownloadUrl(key);
    return info
  }
  async getThumbnailUrl(key: string): Promise<string[]> {
    const info = await this.accessor.backendApi.getThumbnailUrls(key);

    return info.items.map((item)=>item.url)
  }
}
