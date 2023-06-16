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
        items: {
          ...state.items,
          [action.itemId]: ItemReducer(state.items[action.itemId], action),
        },
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
        nextItemId: state.nextArchiveId + action.items.length,
        groups: {
          ...state.groups,
          [action.groupId]: {
            items: {
              ...state.groups[action.groupId].items,
              ...Object.fromEntries(
                action.items.map((item) => [
                  item.id,
                  {
                    file: item.file,
                    marked: false,
                  },
                ])
              ),
            },
          },
        },
      };
    case "CREATE_GROUP":
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.newGroupId]: {
            items: {},
          },
        },
      };
    case "ARCHIVE": {
      const marked: any = {};
      const least: any = {};
      Object.keys(state.groups[action.groupId].items).forEach((key) => {
        if (state.groups[action.groupId].items[key].marked) {
          marked[key] = state.groups[action.groupId].items[key];
        } else {
          least[key] = state.groups[action.groupId].items[key];
        }
      });
      return {
        ...state,
        nextArchiveId: state.nextArchiveId + 1,
        groups: {
          ...state.groups,
          [action.groupId]: {
            items: least,
          },
          [state.nextArchiveId.toString()]: {
            items: marked,
          },
        },
      };
    }
    default:
      return state;
  }
};

export default StatusReducer;
