import AWS from "aws-sdk";
import { uuidv7 } from "uuidv7";
import { BufferUploader } from "./Uploader";

type S3DirectAccessSettings = {
  bucket: string;
  config: AWS.S3.ClientConfiguration;
};

export default class S3DirectAccessBufferUploader implements BufferUploader {
  private readonly s3: AWS.S3;
  private readonly bucket: string;

  constructor(settings: S3DirectAccessSettings) {
    this.s3 = new AWS.S3(settings.config);
    this.bucket = settings.bucket;
  }
  async upload(zipData: ArrayBuffer): Promise<string> {
    const objId = uuidv7();
    const s3Key = `${objId}.zip`;

    const uploadParams = {
      Bucket: this.bucket,
      Key: s3Key,
      Body: zipData,
      ContentType: "application/zip",
      ACL: "public-read",
    };

    await this.s3.upload(uploadParams).promise();

    return objId;
  }
}
