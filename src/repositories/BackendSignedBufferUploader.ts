import { BufferUploader } from "./Uploader";
import axios, { AxiosError } from "axios";

const StringifyGeneralError = (err: Error): string => {
  if (err instanceof AxiosError) {
    if (err.response) {
      return `AxiosError with response status ${err.response.status}, body is "${err.response.data}"`;
    } if (err.request) {
      return `AxiosError with message: ${JSON.stringify(err.toJSON())}`;
    }
  }
  return err.toString();
};

class ApplicationError extends Error {
  constructor(public readonly message: string, public readonly reason?: any) {
    super(message);
  }
  toString(): string {
    return this.toStrings().join()
  }
  toStrings(): string[] {
    let messages: string[] = [];
    if (this.reason) {
      if (this.reason instanceof ApplicationError) {
        messages = this.reason.toStrings();
      } else {
        messages = [StringifyGeneralError(this.reason)];
      }
    }

    messages.unshift(this.message)

    return messages;
  }
}

type Transaction = {
  id: number;
  url: string;
}

class BackendApi {
  constructor(private readonly endpointUrl: string){}
  async createTransaction(): Promise<Transaction> {
    try {
      const response = await axios.post(`${this.endpointUrl}/storage/transactions`, {}, {
        headers: {
          'Content-Type': 'application/json'
        }
      });
 
      return response.data as Transaction
    } catch(err) {
      throw new ApplicationError(`failed to createTransaction`, err)
    }
  }
}

class S3SignedAccessor {
  async put(url: string, buffer: ArrayBuffer): Promise<void>{
    try {
      // Error parsing the X-Amz-Credential parameter; the region 'us-east-1' is wrong; expecting 'ap-northeast-1'
      const response = await axios.put(url, buffer, {
        headers: {
          'Content-Type': 'application/zip'
        }
      });
      console.log(response);
    } catch(err) {
      throw new ApplicationError(`failed to put data to ${url}`, err)
    }
  }
}

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
