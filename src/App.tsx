import { useCallback, useReducer, useState } from "react";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import "./App.css";

const SelectionStatusReducer = (
  state: SelectionStatus[],
  action: SelectionStatusAction
): SelectionStatus[] => {
  switch (action.type) {
    case "SET_MARKED":
      return state.map((v, i) =>
        i === action.index
          ? {
              ...state[action.index],
              marked: action.value,
            }
          : v
      );
    case "LOADED":
      return [
        ...state,
        ...action.files.map((f) => ({ file: f, marked: false })),
      ];
    case "ARCHIVE":
      return state.map((s) =>
        s.marked ? { ...s, status: false, archiveiId: action.archiveId } : s
      );
    default:
      return state;
  }
};

function App() {
  const [currentArchiveId, setCurrentArchiveId] = useState<number>(0);
  const [selectionStatus, dispatch] = useReducer(SelectionStatusReducer, []);

  const onSelect = useCallback(
    (index: number) => {
      dispatch({
        type: "SET_MARKED",
        index: index,
        value: !selectionStatus[index].marked,
      });
    },
    [selectionStatus]
  );

  return (
    <div>
      <div className="gallery-view">
        {selectionStatus.length > 0 ? (
          <>
            {selectionStatus.map((status, index) =>
              status.archiveiId === undefined ? (
                <div key={index} onClick={() => onSelect(index)}>
                  <ImagePreview file={status.file} marked={status.marked} />
                </div>
              ) : null
            )}
          </>
        ) : (
          <FileLoader
            onLoaded={(files) => dispatch({ type: "LOADED", files })}
          />
        )}
      </div>
      <div>
        <button
          onClick={() => {
            dispatch({ type: "ARCHIVE", archiveId: currentArchiveId });
            setCurrentArchiveId(currentArchiveId + 1);
          }}
        >
          Archive
        </button>
      </div>
    </div>
  );
}

export default App;
