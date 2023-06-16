const ItemReducer = (state: Item, action: Action): Item => {
  switch (action.type) {
    case "MARK_ITEM":
      return {
        ...state,
        marked: action.value,
      };
    default:
      return state;
  }
};

const StatusGroupReducer = (state: Group, action: Action): Group => {
  switch (action.type) {
    case "MARK_ITEM":
      return {
        items: state.items.map((v, i) =>
          i === action.index ? ItemReducer(v, action) : v
        ),
      };
    case "LOAD":
    default:
      return state;
  }
};

const StatusReducer = (state: Status, action: Action): Status => {
  switch (action.type) {
    case "MARK_ITEM":
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.groupId]: StatusGroupReducer(
            state.groups[action.groupId],
            action
          ),
        },
      };
    case "LOAD":
      return {
        ...state,
        nextItemId: state.nextArchiveId + action.files.length,
        groups: {
          ...state.groups,
          [action.groupId]: {
            items: [
              ...state.groups[action.groupId].items,
              ...action.files.map((f, i) => ({
                id: state.nextArchiveId + i,
                file: f,
                marked: false,
              })),
            ],
          },
        },
      };
    case "CREATE_GROUP":
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.newGroupId]: {
            items: [],
          },
        }
      };
    case "ARCHIVE":
      return {
        ...state,
        nextArchiveId: state.nextArchiveId + 1,
        groups: {
          ...state.groups,
          [state.nextArchiveId.toString()]: {
            items: state.groups[action.groupId].items.filter((s) => s.marked),
          },
          [action.groupId]: {
            items: state.groups[action.groupId].items.filter((s) => !s.marked),
          },
        },
      };
    default:
      return state;
  }
};

export default StatusReducer;
