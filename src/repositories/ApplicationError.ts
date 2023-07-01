import { AxiosError } from "axios";

export const StringifyGeneralError = (err: Error): string => {
  if (err instanceof AxiosError) {
    if (err.response) {
      return `AxiosError with response status ${err.response.status}, body is "${err.response.data}"`;
    } if (err.request) {
      return `AxiosError with message: ${JSON.stringify(err.toJSON())}`;
    }
  }
  return err.toString();
};

export class ApplicationError extends Error {
  constructor(public readonly message: string, public readonly reason?: any) {
    super(message);
  }
  toString(): string {
    return this.toStrings().join();
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

    messages.unshift(this.message);

    return messages;
  }
}
