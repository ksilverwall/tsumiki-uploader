type LoadAction = {
  type: "SET_GROUP_ITEMS";
  groupId: GroupId;
  items: {[key: ItemId]: Item}
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
  groupId: GroupId;
};

type UploadManyAction = {
  type: "UPLOAD_MANY";
  groupIds: GroupId[];
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
  | UploadManyAction
  | UploadCompleteAction;
