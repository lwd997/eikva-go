import { DocListItem } from "../DocListItem/DocListItem";
import Button from "../universal/Button/Button";
import "./DocList.css";
export const DocList = () => {
  return (
    <div className="d-flex flex-column row-gap-2">
      <Button>Upload</Button>
      <div className="doc-list">
        {[0, 1, 2, 3, 4, 5, 6, 7, 8, 9].map((_, index) => {
          return <DocListItem key={index} />;
        })}
      </div>
    </div>
  );
};
