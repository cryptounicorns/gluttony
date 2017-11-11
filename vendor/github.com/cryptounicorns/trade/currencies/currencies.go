package currencies

type Currencies []Currency

func (cs Currencies) MapByName() map[string]Currency {
	var (
		res = make(
			map[string]Currency,
			len(cs),
		)
	)

	for _, v := range cs {
		res[v.Name] = v
	}

	return res
}

func (cs Currencies) MapBySymbol() map[Symbol]Currency {
	var (
		res = make(
			map[Symbol]Currency,
			len(cs),
		)
	)

	for _, v := range cs {
		res[v.Symbol] = v
	}

	return res
}

func NewFromJSON(buf []byte) (Currencies, error) {
	var (
		res = Currencies{}
		err error
	)

	err = jsonFormat.Unmarshal(buf, &res)

	return res, err
}
