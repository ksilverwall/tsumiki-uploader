import { useCallback, useReducer } from "react";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import "./App.css";
import reducer from "./reducer";

function App() {
  const [status, dispatch] = useReducer(reducer, {
    nextArchiveId: 0,
    nextItemId: 0,
    items: [],
  });

  const onSelect = useCallback(
    (index: number) => {
      dispatch({
        type: "MARK_ITEM",
        index: index,
        value: !status.items[index].marked,
      });
    },
    [status]
  );

  return (
    <div>
      <div className="gallery-view">
        {status.items.length > 0 ? (
          <>
            {status.items
              .filter((item) => item.archiveiId === undefined)
              .map((item, idx) => (
                <div key={idx} onClick={() => onSelect(item.id)}>
                  <ImagePreview file={item.file} marked={item.marked} />
                </div>
              ))}
          </>
        ) : (
          <FileLoader onLoaded={(files) => dispatch({ type: "LOAD", files })} />
        )}
      </div>
      <div>
        <button
          onClick={() => {
            dispatch({ type: "ARCHIVE" });
          }}
        >
          Archive
        </button>
      </div>
      <div>
        {[
          ...new Set(
            status.items
              .map((item) => item.archiveiId)
              .filter((v) => v !== undefined)
          ),
        ].sort().map((v) => (
          <button>{v}</button>
        ))}
      </div>
    </div>
  );
}

export default App;
