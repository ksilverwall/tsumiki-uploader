type SelectionStatusSetMarkedAction = {
  type: "SET_MARKED";
  index: number;
  value: boolean;
};

type SelectionStatusLoadedAction = {
  type: "LOADED";
  files: File[];
};

type SelectionStatusArchiveAction = {
  type: "ARCHIVE";
  archiveId: number;
};

type SelectionStatusAction =
  | SelectionStatusSetMarkedAction
  | SelectionStatusLoadedAction
  | SelectionStatusArchiveAction;
