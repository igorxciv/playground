import { Joke } from "@prisma/client";
import { json, Link, LoaderFunction, useLoaderData } from "remix";
import { db } from "~/utils/db.server";

type LoaderData = {
    randomJoke: Joke
}

export const loader: LoaderFunction = async () => {
    const count = await db.joke.count();
    const randomRowNumber = Math.floor(Math.random() * count);
    const [randomJoke] = await db.joke.findMany({
        take: 1,
        skip: randomRowNumber,
    })

    const data: LoaderData = { randomJoke }

    return json(data);
}

export default function JokesIndex() {
    const { randomJoke } = useLoaderData<LoaderData>()

    return (
        <div>
            <p>Heres random joke:</p>
            <p>
                {randomJoke.content}
            </p>
            <Link to={randomJoke.id}>
                "{randomJoke.name}" Permalink
            </Link>
        </div>
    )
}

export function ErrorBoundary() {
    return (
      <div className="error-container">
        I did a whoopsies.
      </div>
    );
  }