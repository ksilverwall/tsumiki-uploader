import { DefaultApi, TransactionStatus } from "../../gen";
import { BackendInterface, DownloadInfo } from "../interface";
import { ApplicationError } from "./ApplicationError";
import { S3SignedAccessor } from "./S3SignedAccessor";

type Accessor = {
  backendApi: DefaultApi;
  s3Accessor: S3SignedAccessor;
}

export default class Backend implements BackendInterface {
  constructor(private readonly accessor: Accessor) { }
  async upload(zipData: ArrayBuffer): Promise<string> {
    try {
      const res1 = await this.accessor.backendApi.createTransaction();
      await this.accessor.s3Accessor.put(res1.data.url, zipData);
      await this.accessor.backendApi.updateTransaction(res1.data.id, {status: TransactionStatus.Uploaded});
      return res1.data.id;
    } catch (err) {
      throw new ApplicationError(`failed to createTransaction`, err);
    }
  }
  async download(key: string): Promise<DownloadInfo> {
    try {
      const res = await this.accessor.backendApi.getFileUrl(key);
      return res.data;
    } catch (err) {
      throw new ApplicationError(`failed to getDownloadUrl`, err);
    }
  }
  async getThumbnailUrl(key: string): Promise<string[]> {
    try {
      const res = await this.accessor.backendApi.getFileThumbnailUrls(key);

      return res.data.items.map((item) => item.url)
    } catch (err) {
      throw new ApplicationError(`failed to getDownloadUrl`, err);
    }
  }
}
