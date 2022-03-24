package providers

/*
// Основная функция синхронизации данных из VSPHERE
func SyncData(dcs []string) {
	// SYNC TAG
	syncTagExist := SyncTagCheck()
	if !syncTagExist {
		SyncTagCreate()
	}

	ClusterGroupSync(dcs)
}

// ClusterGroupSync
// Функция синхронизации ClusterGroup(DataCenters)
// Принимает в себя массив dcs: новых CG из VSphere и eDcs: массив существующих CG
// TODO: Скорее всего эту функцию нужно будет вынести в другой пакет: providers
func ClusterGroupSync(dcs []string) {
	// Исключить из синхронизации CG из конфига (Exclude)
	exArr := excludeFilter(dcs, cfg.Settings.DataCenters.Exclude)

	// TODO: Включить в синхронизацию CG из конфига (Include)
	// FIXME: inArr := include(exArr, dcs)

	existCgs := ClusterGroupList()
	var cg []string
	for _, o := range existCgs {
		cg = append(cg, *o.Name)
	}

	// Сравнить 2 массива (exArr и eDcs) на выходе должны получить значения которые не имеются во 2м массиве
	addDcs, deleteDcs := diffData(exArr, cg)

	// Create ClusterGroup
	for _, s := range addDcs {
		ClusterGroupCreate(s, s)
	}
	for _, s1 := range deleteDcs {
		for _, s2 := range existCgs {
			if s1 == *s2.Name {
				ClusterGroupDelete(s2.ID)
			}
		}
	}
	fmt.Println("SYNCED")
	// FIXME: ClusterCreate("TEST", 4)
	fmt.Println(ClusterTypeCheck("VMware ESXi"))
	fmt.Println(VmList())
}
*/
