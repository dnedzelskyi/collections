const API = {
  getData: async (searchText = '') => {
    const url = new URL(`${document.URL}/api/data`);
    searchText && url.searchParams.append('s', searchText);
    return (await fetch(url)).json();
  },
};

document.addEventListener('DOMContentLoaded', init);

async function init() {
  document
    .querySelector('#search')
    ?.addEventListener('input', debounce(handleSearch, 1000));
  handleSearch();
}

async function handleSearch(event) {
  let items = (await API.getData(event?.target?.value)) ?? [];
  renderSearchResults(items);
}

function renderSearchResults(items = []) {
  const searchResultContainer = document.querySelector(
    '.search-result-container',
  );
  const fragment = document.createDocumentFragment();
  const cardTemplate = document.querySelector(`#card`);

  for (let data of items) {
    let card = cardTemplate.content.cloneNode(true).querySelector('article');

    card.dataset.id = data?.id;
    card.dataset.createdAt = data?.createdAt;

    let [openDetailsPopoverButton, detailsPopover] = [
      card.querySelector('.card-btn'),
      card.querySelector('.card-details-popover'),
    ];
    openDetailsPopoverButton.popoverTargetElement = detailsPopover;

    for (let [key, value] of Object.entries(data)) {
      const element = card.querySelector(`[data-field="${key}"]`);
      element && (element.textContent = value);
    }

    fragment.appendChild(card);
  }

  searchResultContainer.replaceChildren(fragment);
}

function debounce(fn, delay = 3000) {
  let timeoutId = undefined;
  return (...args) => {
    timeoutId && clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn?.(...args), delay);
  };
}
