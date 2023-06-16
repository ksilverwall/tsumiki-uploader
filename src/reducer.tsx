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

const StatusReducer = (state: Status, action: Action): Status => {
  switch (action.type) {
    case "MARK_ITEM":
      return {
        ...state,
        items: state.items.map((v, i) =>
          i === action.index ? ItemReducer(v, action) : v
        ),
      };
    case "LOAD":
      return {
        ...state,
        nextItemId: state.nextArchiveId + action.files.length,
        items: [
          ...state.items,
          ...action.files.map((f, i) => ({
            id: state.nextArchiveId + i,
            file: f,
            marked: false,
          })),
        ],
      };
    case "ARCHIVE":
      return {
        ...state,
        nextArchiveId: state.nextArchiveId + 1,
        items: state.items.map((s) =>
          s.marked
            ? {
                ...s,
                marked: false,
                archiveId: state.nextArchiveId,
              }
            : s
        ),
      };
    default:
      return state;
  }
};

export default StatusReducer;
