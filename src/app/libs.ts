import JSZip from "jszip";
import { uuidv7 } from "uuidv7";

export async function ArchiveFiles(files: File[]): Promise<ArrayBuffer> {
    const zip = new JSZip();

    files.forEach((file) => {
        zip.file(file.name, file);
    });

    return await zip.generateAsync({ type: "arraybuffer" });
}

export function GenerateId<T extends string>(): T {
  const newId = uuidv7();
  return newId as T;
}
