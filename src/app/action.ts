import { GroupId, Item, ItemId } from "./models";
import { ApplicationError } from "./repositories/ApplicationError";

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

type UploadFailedAction = {
  type: "UPLOAD_FAILED";
  groupId: GroupId;
  error: ApplicationError;
};

export type Action =
  | LoadAction
  | CreateGroupAction
  | UploadAction
  | UploadManyAction
  | UploadCompleteAction
  | UploadFailedAction;
