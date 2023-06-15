type SelectionStatusSetMarkedAction = {
  type: "SET_MARKED";
  index: number;
  value: boolean;
};

type SelectionStatusLoadedAction = {
  type: "LOADED";
  files: File[];
};

type SelectionStatusAction =
  | SelectionStatusSetMarkedAction
  | SelectionStatusLoadedAction;
