# NetBox Sync WIP
Program for synchronizing data from virtualization systems

 - [ ] VSphere (WIP)
 - [ ] Nutanix (https://github.com/tecbiz-ch/nutanix-go-sdk)
 - [ ] HyperV?
 - [ ] etc


Принцип работы:
1. Получить все данные из NetBox
  - Цепочка данных в NB: 
    
    ClusterGroup ->
                    Cluster -> VMs
    ClusterType  -> 
    
    
    ClusterGroup
      Нельзя добавить тег. (Поиск решения)
      VSphere: DataCenter

      - ClusterGroupCreate  - Создать новую кластер-группу
      - ClusterGroupDelete  - Удалить существующую кластер-группу
      - ClusterGroupList    - Получить существующие кластер-группы

    ClusterType
      Нельзя добавить тег. (Поиск решения)
      VSphere: VMware ESXi (Хардкодед значение)

      - ClusterTypeCreate  - Создать новый тип-кластера
      - ClusterTypeDelete  - Удалить существующий тип-кластера
      - ClusterTypeList    - Получить существующиие типы-кластера
      - ClusterTypeCheck   - Проверить существует ли указанный тип кластера (Поиск по имени)

    Cluster
      VSphere: Cluster
      При создании необходимы данные:
        Name, ClusterTypeID, ClusterGroupID, Tag

      - ClusterCreate - Создать новый кластер
      - ClusterList   - Получить список кластеров по тегу

    VMs
      При создании необходимы данные:
        Name, ClusterID, Tag and etc

      - VmList   - Получить список виртуальных машин по тегу

2. Получить все данные из источников (VSphere)
3. Обработать данные, получить что нужно добавить, обновить и удалить
4. Внести изменения


TODO:

- ClusterGroups (Datacenters)
  - ISSIE:
    - Невозможно добавить теги, не позволяет библиотека (Написал запрос)
  + Возможность исключить определенные датацентры
  - Возможность включить определенные датацентры
  + Если датацентр существует -> пропустить

- VMs
  - Виртуальные машины в Netbox должны иметь уникальные имена в кластере. При дублях выводить ошибку с добавлением ВМ (Название ВМ, кластер)

- SyncTag
  - При инициализации проверяем создан ли тег, по которому мы будем определять что под нашим конктролем.
    + Функция создание тега
    + Функция проверки тега по названию, возвращает тру ор фалс, если фалс то создать тег
    + Сохранить ИД тега который потом будем прикреплять к итемам в глобальной переменной

- Include / Exclude Filters
  - Как должен работать фильтр Include относительно фильтра Exclude?
  - Сначала прогоняем слайс через Include Regex Filter, и потом через Exclude Regex Filter?