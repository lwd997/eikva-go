import Icon from "../components/universal/Icon/Icon";

export const Home = () => {
    return (
        <div className="card group height-100 width-100 display-flex justify-content-center align-items-center">
            <div className="width-50">
                <div className="logo">
                    <img src="/TestCraft.svg" />
                </div>

                <h2>Начало работы</h2>
                <p>
                    Для начала работы создайте новую группу или выберите существующую
                </p>

                <h2>Создание тест-кесов</h2>
                <p>
                    Для генерации тест-кейсов, откройте рабочую панель, находясь в группе. <br />
                    Прикрепите файлы или введите текст в поле ввода, укажите количество кейсов
                    и нажмите "Начать генерацию"
                </p>

                <p>
                    Так же кейсы можно создавать самостоятельно, нажав на "Добавить новый тест-кейс"
                </p>

                <p>
                    Вы можете редактировать и удалять кейсы, созданные вами.
                </p>

                <h2>Экспорт тест-кейсов</h2>
                <p>
                    Для экспорта группы тест-кейсов, нажмите на кнопку <Icon name="download" />
                    рядом с соответствующим пунктом в списке групп.
                </p>
                <p>
                    Выберите формат файла: <br /><br />
                    <Icon name="table" /> - Excel (xslx) <br />
                    <Icon name="data_object" /> - Zephyr (json)
                </p>
            </div>
        </div>
    );
};
