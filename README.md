# avito_intern
Запуск проекта : make run (билд и поднятие контейнеров)<br/>
Тесты: make test<br/><br/>
Документация описана в swagger по юрлу /docs .Стоит лишь добавить ,что модель ответа одинаковая - {status:"status","error":"optional","response":"optional"}.В response возвращаемые поля
.Также существует несколько ошибок :validation error (ошибка валидации),server error(ошибка бд или сервиса),already exists. <br/>Юрлы разделены по тэгам:User - операции с пользователем,Segment -  операции с сегментами ,History  - первое доп.задание  <br/><br/>
Выполнены все доп задания,однако возникли некоторые трудности:<br/>
- Не совсем понял , что нужно возвращать в первом задании(история добавления/удаления). Я решил ,что просто csv файл ,а не ссылка
- В третьем задании я решил,что существующие пользователи получают сегменты


