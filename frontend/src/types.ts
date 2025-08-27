export type Artist = {
	id: number;
	name: string;
	image: string;
	customImage: string;
	members: string[];
	creationDate: number;
	firstAlbum: string;
	locations: string;
	concertDates: string;
	relations: string;
};

export type Location = {
	name: string;
	lat: string;
	lon: string;
	boundingbox: [string, string, string, string];
	type: string;
	dates: string[];
};
