import { BufferUploader } from "./Uploader";
import axios from "axios";

type Transaction = {
  id: number;
  url: string;
}

class BackendApi {
  constructor(private readonly endpointUrl: string){}
  async createTransaction(): Promise<Transaction> {
      // FIXME: blocked by CORS policy
      const response = await axios.post(`${this.endpointUrl}/storage/transactions`, {}, {
        headers: {
          'Content-Type': 'application/json'
        }
      });
 
      return response.data as Transaction
  }
}

export default class BackendSignedBufferUploader implements BufferUploader {
  constructor() {}
  async upload(zipData: ArrayBuffer): Promise<string> {
    try {
      const endPoint = import.meta.env.VITE_BACKEND_API_ENDPOINT;
      const backendApi = new BackendApi(endPoint);
      const t = await backendApi.createTransaction();

      console.log(t);
  
      const response2 = await axios.put(t.url, zipData, {
        headers: {
          'Content-Type': 'application/zip'
        }
      });
      console.log(response2);
    } catch (err) {
      console.error('Error sending PUTObject request:', err);
    }
  
    return "";
  }
}
