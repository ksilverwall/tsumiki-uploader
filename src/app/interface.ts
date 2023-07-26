export type DownloadInfo = {
  name: string;
  url: string;
};

export interface BackendInterface {
  upload(zipData: ArrayBuffer): Promise<string>;
  download(key: string): Promise<DownloadInfo>;
  getThumbnailUrl(key: string): Promise<string[] | undefined>;
}
