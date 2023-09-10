const growers = document.querySelectorAll(".auto-grow-textarea-wrapper");

growers.forEach((grower) => {
  const textarea = grower.querySelector("textarea");
  textarea.addEventListener("input", () => {
    grower.dataset.replicatedValue = textarea.value;
  });
});