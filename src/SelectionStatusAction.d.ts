type MarkItemAction = {
  type: "MARK_ITEM";
  index: number;
  value: boolean;
};

type LoadAction = {
  type: "LOAD";
  files: File[];
};

type ArchiveAction = {
  type: "ARCHIVE";
};

type Action =
  | MarkItemAction
  | LoadAction
  | ArchiveAction;
