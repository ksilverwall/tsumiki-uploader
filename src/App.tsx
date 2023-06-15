import { useCallback, useReducer } from "react";
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
    default:
      return state;
  }
};

function App() {
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
    <div className="gallery-view">
      {selectionStatus.length > 0 ? (
        <>
          {selectionStatus.map((status, index) => (
            <div key={index} onClick={() => onSelect(index)}>
              <ImagePreview file={status.file} marked={status.marked} />
            </div>
          ))}
        </>
      ) : (
        <FileLoader onLoaded={(files) => dispatch({ type: "LOADED", files })} />
      )}
    </div>
  );
}

export default App;
