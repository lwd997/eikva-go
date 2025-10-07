import Button from "../universal/Button/Button";
import "./UploadDetails.css";
export const UploadDetails = () => {
  return (
    <div className="upload-details">
      <textarea
        className="form-control w-100 d-block"
        style={{
          padding: "1rem",
          height: "6.75rem",
        }}
      />
      <div className="actions-group">
        <Button className="action">Generate</Button>
        {[0, 1, 2, 3].map((_, index) => {
          return (
            <Button className="action" key={index}>
              Action 0{index + 2}
            </Button>
          );
        })}
      </div>
    </div>
  );
};
