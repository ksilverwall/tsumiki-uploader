type Item = {
  id: number;
  file: File;
  marked: boolean;
  archiveiId?: number;
};

type Status = {
  nextArchiveId: number;
  nextItemId: number;
  items: Item[];
};
