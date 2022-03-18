import { Joke } from "@prisma/client";
import { ActionFunction, Form, json, Link, LoaderFunction, MetaFunction, redirect, useCatch, useLoaderData, useParams } from "remix"
import { JokeDisplay } from "~/components/joke";

import { db } from "~/utils/db.server"
import { getUserId, requireUserId } from "~/utils/session.server"

type LoaderData = {
    joke: Joke,
    isOwner: boolean;
}

export const meta: MetaFunction = ({ data }: {data: LoaderData}) => {
    if (!data) {
        return {
            title: 'No joke',
            description: 'No joke found',
        }
    }
    return {
        title: `${data.joke.name} joke`,
        description: `Enjoy the ${data.joke.name} joke and much more`,
    }
}

export const loader: LoaderFunction = async ({ params, request }) => {
    const { jokeId } = params;
    const userId = await getUserId(request);

    const joke = await db.joke.findUnique({
        where: { id: jokeId },
    });
    if (!joke) throw new Response('What a joke! Not found', { status: 404 });

    const data: LoaderData = {
        joke,
        isOwner: userId === joke.jokesterId
    }

    return json(data);
}

export const action: ActionFunction = async ({ request, params }) => {
    const form = await request.formData()
    if (form.get('_method') !== 'delete') {
        throw new Response(`The method ${form.get('_method')} is not supported`, { status: 405 });
    }
    const userId = await requireUserId(request);
    const joke = await db.joke.findUnique({ where: { id: params.jokeId } })
    if (!joke) {
        throw new Response('Cant delete what does not exist', { status: 404 });
    }

    if (joke.jokesterId !== userId) {
        throw new Response('Pssh, nice try. Thats not your joke', { status: 403 });
    }

    await db.joke.delete({ where: { id: params.jokeId } })
    return redirect('/jokes');
}

export default function JokeRoute() {
    const { joke, isOwner } = useLoaderData<LoaderData>()
    return (
        <JokeDisplay joke={joke} isOwner={isOwner} />
    )
}

export function ErrorBoundary() {
    const { jokeId } = useParams()

    return (
        <div className="error-container">{`There was an error loading joke by the id ${jokeId}. Sorry.`}</div>
    )
}

export function CatchBoundary() {
    const caught = useCatch()
    const params = useParams()
    switch (caught.status) {
        case 404: {
            return (
                <div className="error-container">Huh? What the heck is "{params.jokeId}"?</div>
            )
        }
        case 403: {
            return (
                <div className="error-container">
                    Sorry, but {params.jokeId} is not your joke.
                </div>
            )
        }
        case 405: {
            return (
                <div className="error-container">
                    What you're trying to do is not allowed.
                </div>
            )
        }
    }
    throw new Error(`Unhandled error: ${caught.status}`)
}