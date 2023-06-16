type GroupId = string;
type ItemId = string;

type Item = {
  file: File;
  marked: boolean;
};

type Group = {
  items: { [key: ItemId]: Item };
};

type Status = {
  nextArchiveId: number;
  nextItemId: number;
  groups: { [key: GroupId]: Group };
};
