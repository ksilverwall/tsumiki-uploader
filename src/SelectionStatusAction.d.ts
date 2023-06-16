type LoadAction = {
  type: "LOAD";
  groupId: GroupId;
  items: { id: ItemId; file: File }[];
};

type CreateGroupAction = {
  type: "CREATE_GROUP";
  newGroupId: GroupId;
  source?: {
    groupId: GroupId;
    selected: ItemId[];
  };
};

type UploadAction = {
  type: "UPLOAD";
};

type UploadCompleteAction = {
  type: "UPLOAD_COMPLETE";
  groupId: GroupId;
  key: string;
};

type Action =
  | LoadAction
  | CreateGroupAction
  | UploadAction
  | UploadCompleteAction;
