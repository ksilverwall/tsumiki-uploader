type GroupId = string;
type ItemId = string;

type Item = {
  id: number;
  file: File;
  marked: boolean;
};

type Group = {
  items: Item[];
};

type Status = {
  nextArchiveId: number;
  nextItemId: number;
  groups: { [key: GroupId]: Group };
};
