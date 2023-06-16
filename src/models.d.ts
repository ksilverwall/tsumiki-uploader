type Item = {
  id: number;
  file: File;
  marked: boolean;
  archiveId: number;
};

type Status = {
  nextArchiveId: number;
  nextItemId: number;
  items: Item[];
};
