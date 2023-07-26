import { AxiosError } from "axios";
import { ClientErrorCode, DefaultApi, GetFileUrl400Response, GetFileUrl500Response, TransactionStatus } from "../../gen";
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
      const res2 = await this.accessor.backendApi.updateTransaction(res1.data.id, { status: TransactionStatus.Uploaded });

      return res2.data.file_id;
    } catch (err) {
      throw new ApplicationError(`failed to upload zip file`, err);
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
  async getThumbnailUrl(key: string): Promise<string[] | undefined> {
    try {
      const res = await this.accessor.backendApi.getFileThumbnailUrls(key);

      return res.data.items.map((item) => item.url)
    } catch (err) {
      if (err instanceof AxiosError && err.response) {
        if (Math.floor(err.response.status / 100) === 4) {
          const d = err.response.data as GetFileUrl400Response
          switch (d.code) {
          case ClientErrorCode.ThumbnailNotFound:
            return undefined;
          }
        }
      }

      throw new ApplicationError(`failed to getDownloadUrl`, err);
    }
  }
}
